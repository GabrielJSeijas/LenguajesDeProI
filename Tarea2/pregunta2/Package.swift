// swift-tools-version: 6.2
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "pregunta2",
    products: [
        // Products define the executables and libraries a package produces, making them visible to other packages.
        .library(
            name: "pregunta2",
            targets: ["pregunta2"]
        ),
    ],
    targets: [
    // 1. Target de Librería (Tu lógica de EVAL y MOSTRAR)
    .target(
        name: "pregunta2",
        dependencies: []),

    // 2. Target Ejecutable (Tu consola interactiva)
    .executableTarget(
        name: "App",
        dependencies: ["pregunta2"]), // Depende de tu librería

    // 3. Target de Pruebas
    .testTarget(
        name: "pregunta2Tests",
        dependencies: ["pregunta2"]),
]
)
