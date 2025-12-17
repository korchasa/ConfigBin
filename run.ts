#!/usr/bin/env -S deno run --allow-read --allow-write --allow-run --allow-env

/**
 * ConfigBin run script
 * Automation script for Go project management
 */

interface Command {
  name: string;
  description: string;
  handler: (args: string[]) => Promise<number>;
}

const commands: Map<string, Command> = new Map();

// Utility functions
async function exec(
  cmd: string[],
  options: { cwd?: string; env?: Record<string, string> } = {},
): Promise<{ success: boolean; code: number; output: string }> {
  const process = new Deno.Command(cmd[0], {
    args: cmd.slice(1),
    cwd: options.cwd,
    env: options.env,
    stdout: "piped",
    stderr: "piped",
  });

  const { code, stdout, stderr } = await process.output();
  const output = new TextDecoder().decode(stdout) +
    new TextDecoder().decode(stderr);

  return {
    success: code === 0,
    code,
    output,
  };
}

function log(message: string): void {
  console.log(`[run.ts] ${message}`);
}

function error(message: string): void {
  console.error(`[run.ts] ERROR: ${message}`);
}

function checkCommand(cmd: string): Promise<boolean> {
  return exec(["which", cmd]).then((r) => r.success);
}

// Commands
async function cmdInit(_args: string[]): Promise<number> {
  log("Installing dependencies...");
  const steps = [
    { cmd: ["go", "mod", "download"], desc: "Downloading modules" },
    { cmd: ["go", "mod", "vendor"], desc: "Vendoring dependencies" },
  ];

  for (const step of steps) {
    log(step.desc);
    const result = await exec(step.cmd);
    if (!result.success) {
      error(`${step.desc} failed: ${result.output}`);
      return result.code;
    }
  }

  log("Dependencies installed successfully");
  return 0;
}

async function cmdTestOne(args: string[]): Promise<number> {
  if (args.length === 0) {
    error("test-one requires a test path argument");
    console.log("Usage: ./run.ts test-one <test_path>");
    console.log("Example: ./run.ts test-one ./pkg/encryptor/aes");
    console.log("Example: ./run.ts test-one configBin/pkg/encryptor/aes");
    return 1;
  }

  const testPath = args[0];
  log(`Running test: ${testPath}`);

  // Try with ./ prefix if not present
  let path = testPath;
  if (!testPath.startsWith("./") && !testPath.startsWith("configBin/")) {
    // Check if it's a relative path that needs ./ prefix
    try {
      const stat = await Deno.stat(testPath).catch(() => null);
      if (!stat) {
        // Try with ./ prefix
        const altPath = `./${testPath}`;
        const altStat = await Deno.stat(altPath).catch(() => null);
        if (altStat) {
          path = altPath;
        } else {
          // Try as module path
          path = `configBin/${testPath}`;
        }
      }
    } catch {
      // If path doesn't exist, try module path
      path = `configBin/${testPath}`;
    }
  }

  const result = await exec(["go", "test", "-v", path]);
  if (!result.success) {
    error(`Test failed: ${result.output}`);
    return result.code;
  }

  console.log(result.output);
  return 0;
}

async function cmdCleanup(_args: string[]): Promise<number> {
  log("Cleaning project...");

  const itemsToRemove = [
    "app",
    "bin/run/app",
    "bin/run/configBin",
    "*.test",
    "coverage.out",
    "var",
  ];

  for (const item of itemsToRemove) {
    try {
      const stat = await Deno.stat(item).catch(() => null);
      if (stat) {
        if (stat.isDirectory) {
          await Deno.remove(item, { recursive: true });
          log(`Removed directory: ${item}`);
        } else {
          await Deno.remove(item);
          log(`Removed file: ${item}`);
        }
      }
    } catch (e) {
      // Ignore errors for non-existent files
      if (!(e instanceof Deno.errors.NotFound)) {
        error(`Failed to remove ${item}: ${e.message}`);
      }
    }
  }

  // Remove test binaries
  try {
    const entries = await Array.fromAsync(Deno.readDir("."));
    for (const entry of entries) {
      if (entry.name.endsWith(".test") && entry.isFile) {
        await Deno.remove(entry.name);
        log(`Removed test binary: ${entry.name}`);
      }
    }
  } catch (e) {
    // Ignore if directory read fails
  }

  log("Cleanup completed");
  return 0;
}

async function cmdDev(_args: string[]): Promise<number> {
  log("Starting development mode...");

  // Check if gin is available (hot reload tool)
  const hasGin = await checkCommand("gin");
  const listen = Deno.env.get("LISTEN") || "localhost:8080";
  const sqlitePath = Deno.env.get("SQLITE_PATH") || "./var/configbin.db";

  if (hasGin) {
    log("Using gin for hot reload");
    const result = await exec([
      "gin",
      "-i",
      "--port",
      "28080",
      "--appPort",
      "8080",
      "--build",
      "bin/run",
      "--all",
      "--excludeDir",
      "./var",
      "--excludeDir",
      "./vendor",
      "--excludeDir",
      "./.git",
      "run",
    ], {
      env: {
        LISTEN: listen,
        SQLITE_PATH: sqlitePath,
      },
    });

    return result.code;
  } else {
    log("gin not found, running directly");
    log("Note: Install gin for hot reload: go install github.com/codegangsta/gin@latest");
    const result = await exec(["go", "run", "bin/run/main.go"], {
      env: {
        LISTEN: listen,
        SQLITE_PATH: sqlitePath,
      },
    });

    if (!result.success) {
      error(`Development server failed: ${result.output}`);
      return result.code;
    }

    return 0;
  }
}

// Check command stages
async function stageClean(): Promise<number> {
  log("Stage: clean");
  return await cmdCleanup([]);
}

async function stageCompile(): Promise<number> {
  log("Stage: compile");
  const result = await exec(["go", "build", "-o", "app", "bin/run/main.go"]);
  if (!result.success) {
    error(`Compilation failed: ${result.output}`);
    return result.code;
  }
  log("Compilation successful");
  return 0;
}

async function stageCommentScan(): Promise<number> {
  log("Stage: comment-scan");

  const patterns = [
    { pattern: /TODO|FIXME/gi, name: "TODO/FIXME" },
    { pattern: /fmt\.Print|print\(|println\(/gi, name: "Debug prints" },
    { pattern: /nolint|nolint:|#nolint/gi, name: "Linter suppressions" },
  ];

  let found = false;
  const excludeDirs = ["vendor", ".git", "var", "node_modules"];

  async function scanDir(dir: string): Promise<void> {
    try {
      for await (const entry of Deno.readDir(dir)) {
        const fullPath = `${dir}/${entry.name}`;

        if (entry.isDirectory) {
          if (!excludeDirs.includes(entry.name)) {
            await scanDir(fullPath);
          }
        } else if (entry.name.endsWith(".go")) {
          const content = await Deno.readTextFile(fullPath);
          const lines = content.split("\n");

          for (let i = 0; i < lines.length; i++) {
            const line = lines[i];
            for (const { pattern, name } of patterns) {
              if (pattern.test(line)) {
                if (!found) {
                  console.log("\nFound comments/issues:");
                  found = true;
                }
                console.log(`${fullPath}:${i + 1} - ${name}`);
                console.log(`  ${line.trim()}`);
              }
            }
          }
        }
      }
    } catch (e) {
      // Ignore permission errors
    }
  }

  await scanDir(".");
  if (found) {
    log("Comment scan completed (issues found above)");
    return 1;
  } else {
    log("Comment scan completed (no issues found)");
    return 0;
  }
}

async function stageFormat(): Promise<number> {
  log("Stage: format");
  const result = await exec(["go", "fmt", "./..."]);
  if (!result.success) {
    error(`Formatting failed: ${result.output}`);
    return result.code;
  }
  log("Formatting completed");
  return 0;
}

async function stageLint(): Promise<number> {
  log("Stage: lint");

  const hasGolangci = await checkCommand("golangci-lint");
  if (!hasGolangci) {
    error("golangci-lint not found. Install it: https://golangci-lint.run/usage/install/");
    return 1;
  }

  const result = await exec(["golangci-lint", "run", "./..."]);
  if (!result.success) {
    error(`Linting failed: ${result.output}`);
    console.log(result.output);
    return result.code;
  }
  log("Linting completed");
  return 0;
}

async function stageTest(): Promise<number> {
  log("Stage: test");
  const result = await exec(["go", "test", "-v", "./..."]);
  if (!result.success) {
    error(`Tests failed: ${result.output}`);
    console.log(result.output);
    return result.code;
  }
  log("Tests completed");
  return 0;
}

async function stageAnalyze(): Promise<number> {
  log("Stage: analyze");
  const result = await exec(["go", "vet", "./..."]);
  if (!result.success) {
    error(`Analysis failed: ${result.output}`);
    console.log(result.output);
    return result.code;
  }
  log("Analysis completed");
  return 0;
}

async function cmdCheck(args: string[]): Promise<number> {
  const stages: Map<string, () => Promise<number>> = new Map([
    ["clean", stageClean],
    ["compile", stageCompile],
    ["comment-scan", stageCommentScan],
    ["format", stageFormat],
    ["lint", stageLint],
    ["test", stageTest],
    ["analyze", stageAnalyze],
  ]);

  if (args.length > 0) {
    const stageName = args[0];
    const stage = stages.get(stageName);
    if (!stage) {
      error(`Unknown stage: ${stageName}`);
      console.log("Available stages:", Array.from(stages.keys()).join(", "));
      return 1;
    }
    return await stage();
  }

  // Run all stages sequentially
  const stageOrder = [
    "clean",
    "compile",
    "comment-scan",
    "format",
    "lint",
    "test",
    "analyze",
  ];

  log("Running check with all stages...");
  for (const stageName of stageOrder) {
    const stage = stages.get(stageName);
    if (!stage) continue;

    const code = await stage();
    if (code !== 0) {
      error(`Stage '${stageName}' failed`);
      return code;
    }
  }

  log("All check stages completed successfully");
  return 0;
}

// Register commands
commands.set("init", {
  name: "init",
  description: "Install project dependencies",
  handler: cmdInit,
});

commands.set("test-one", {
  name: "test-one",
  description: "Run specific test by path",
  handler: cmdTestOne,
});

commands.set("cleanup", {
  name: "cleanup",
  description: "Clean the project directory",
  handler: cmdCleanup,
});

commands.set("dev", {
  name: "dev",
  description: "Run project in development mode",
  handler: cmdDev,
});

commands.set("check", {
  name: "check",
  description: "Run check stages (clean, compile, comment-scan, format, lint, test, analyze)",
  handler: cmdCheck,
});

// Main
async function main(): Promise<number> {
  const args = Deno.args;

  if (args.length === 0) {
    console.log("ConfigBin run script");
    console.log("\nUsage: ./run.ts <command> [options]");
    console.log("\nCommands:");
    for (const cmd of commands.values()) {
      console.log(`  ${cmd.name.padEnd(15)} ${cmd.description}`);
    }
    console.log("\nCheck stages:");
    console.log("  clean          Clean the project");
    console.log("  compile        Compile the project");
    console.log("  comment-scan   Scan for TODOs, FIXMEs, debug prints, linter suppressions");
    console.log("  format         Format code");
    console.log("  lint           Lint code");
    console.log("  test           Run tests");
    console.log("  analyze        Run static analysis");
    console.log("\nExamples:");
    console.log("  ./run.ts init");
    console.log("  ./run.ts test-one pkg/encryptor/aes");
    console.log("  ./run.ts check");
    console.log("  ./run.ts check lint");
    console.log("  ./run.ts dev");
    return 0;
  }

  const cmdName = args[0];
  const cmdArgs = args.slice(1);

  const cmd = commands.get(cmdName);
  if (!cmd) {
    error(`Unknown command: ${cmdName}`);
    console.log("Run ./run.ts for help");
    return 1;
  }

  return await cmd.handler(cmdArgs);
}

if (import.meta.main) {
  Deno.exit(await main());
}

