export type GrpcRouteEntry = {
    grpcMethod: string;
    resolverName: string;
    methodName: string;
};

type PendingGrpcRouteEntry = {
    grpcMethod: string;
    methodName: string;
    controllerConstructor: object;
};

const pendingGrpcRouteEntries: PendingGrpcRouteEntry[] = [];
const controllerResolverByConstructor = new WeakMap<object, string>();

export function GrpcController(resolverName: string) {
    return function (constructor: object) {
        controllerResolverByConstructor.set(constructor, resolverName);
    };
}

export function GrpcMethod(grpcMethod: string) {
    return function (target: object, propertyKey: string, _descriptor: PropertyDescriptor) {
        pendingGrpcRouteEntries.push({
            grpcMethod,
            methodName: propertyKey,
            controllerConstructor: (target as { constructor: object }).constructor,
        });
    };
}

export function getGrpcRoutes(): GrpcRouteEntry[] {
    return pendingGrpcRouteEntries.map((entry) => {
        const resolverName = controllerResolverByConstructor.get(entry.controllerConstructor);
        if (!resolverName) {
            throw new Error(`@GrpcMethod used without @GrpcController on ${entry.methodName}`);
        }
        return {
            grpcMethod: entry.grpcMethod,
            resolverName,
            methodName: entry.methodName,
        };
    });
}
