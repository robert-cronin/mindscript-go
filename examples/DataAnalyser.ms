agent DataAnalyser {
    goal: "Analyse data and generate reports";
    capabilities: ["Data Analysis", "Report Generation"];
    
    behavior {
        on "new analysis request" {
            var rawData: int = 56 * 8;
            var AnalysedData: float = Analyse(rawData, 2.71);
            log(AnalysedData);
            syscall("mkdir", "analysis-results");
            var report: string = exec("python", "generate_report.py");
            log(report);
        }
    }

    function Analyse(input: int, factor: float): float {
        return input / factor;
    }
}