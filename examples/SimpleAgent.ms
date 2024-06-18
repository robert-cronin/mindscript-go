agent SimpleAgent {
    goal: "Simple agent, only performs the build in functions";
    capabilities: ["Syscalls", "Log"];
    
    behavior {
        on "start" {
            log("Agent started");
            syscall("print", "Hello World");
        }
    }
}