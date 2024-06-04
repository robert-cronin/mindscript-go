agent DataCollector {
    goal: "Collect data from various sources";
    capabilities: ["Data Collection"];

    behavior {
        on "new collection request" {
            var data: string = exec("python", "collect_data.py");
            log(data);
        }
    }
}

agent DataAnalyser {
    goal: "Analyse data and generate reports";
    capabilities: ["Data Analysis", "Report Generation"];

    behavior {
        on "new analysis request" {
            var rawData: int = 56 * 8;
            var analysedData: float = analyse(rawData, 2.71);
            log(analysedData);
            syscall("mkdir", "analysis-results");
            var report: string = exec("python", "generate_report.py");
            log(report);
        }
    }

    function analyse(input: int, factor: float): float {
        return input / factor;
    }
}

agent ReportDistributor {
    goal: "Distribute reports to stakeholders";
    capabilities: ["Report Distribution"];

    behavior {
        on "new distribution request" {
            var report: string = exec("cat", "analysis-results/report.txt");
            syscall("mail", "-s", "New Report", "stakeholder@example.com", "<", report);
        }
    }
}