agent DataProcessor {
    goal: "Process data and perform system operations";
    capabilities: ["Computation", "Syscalls"];
    
    behavior {
        on "new data" {
            var data: int = 42 * 7;
            var result: float = compute(data, 3.14);
            log(result);
            syscall("ls", "-la");
            var scriptOutput: string = exec("python", "script.py");
            log(scriptOutput);
        }
    }

    function compute(input: int, multiplier: float): float {
        return input * multiplier;
    }
}
