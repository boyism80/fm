import { runMigrate } from "./sequelize-migrator";

runMigrate().catch((error: unknown) => {
    console.error(error);
    process.exitCode = 1;
});
