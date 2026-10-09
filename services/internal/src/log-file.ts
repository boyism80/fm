import fs from "fs";
import path from "path";
import { format } from "util";

const UTF8_BOM = "\uFEFF";

function pad(value: number): string {
    return String(value).padStart(2, "0");
}

export class DailyLogFile {
    private currentPath = "";
    private fd: number | null = null;

    constructor(private readonly service: string) {
    }

    write(message: string) {
        const now = new Date();
        const date = `${now.getFullYear()}-${pad(now.getMonth() + 1)}-${pad(now.getDate())}`;
        const time = `${pad(now.getHours())}:${pad(now.getMinutes())}:${pad(now.getSeconds())}`;
        const filePath = path.join("logs", `${date}-${this.service}.log`);
        if (filePath !== this.currentPath) {
            this.open(filePath);
        }

        fs.writeSync(this.fd!, `${date.split("-").join("/")} ${time} ${message}\n`);
    }

    private open(filePath: string) {
        fs.mkdirSync(path.dirname(filePath), { recursive: true });
        const fd = fs.openSync(filePath, "a");
        if (fs.fstatSync(fd).size === 0) {
            fs.writeSync(fd, UTF8_BOM);
        }

        if (this.fd !== null) {
            fs.closeSync(this.fd);
        }
        this.currentPath = filePath;
        this.fd = fd;
    }

    captureConsole() {
        for (const method of ["log", "info", "warn", "error"] as const) {
            const original = console[method].bind(console);
            console[method] = (...args: unknown[]) => {
                original(...args);
                try {
                    this.write(format(...args));
                } catch {
                }
            };
        }
    }
}
