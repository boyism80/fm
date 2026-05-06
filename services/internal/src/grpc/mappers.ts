import { createMapper } from "@automapper/core";
import { pojos } from "@automapper/pojos";

export const grpcMapper = createMapper({
    strategyInitializer: pojos(),
});
