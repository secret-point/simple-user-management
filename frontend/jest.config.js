export default {
  preset: "ts-jest",
  testEnvironment: "node",
  moduleNameMapper: {
    "^@/(.*)$": "<rootDir>/src/$1",
  },
  moduleFileExtensions: ["js", "ts", "json", "vue"],
  transform: {
    "^.+\\.js$": "babel-jest", // Transform JS with Babel
    "^.+\\.ts$": "ts-jest", // Transform TypeScript with ts-jest
    "^.+\\.vue$": "@vue/vue3-jest",
    "^.+\\js$": "babel-jest",
  },
  // preset: "ts-jest",
  testEnvironment: "jsdom",
  // transform: {
  //   "^.+\\.vue$": "@vue/vue3-jest",
  //   "^.+\\js$": "babel-jest",
  // },
  testRegex: "(/__tests__/.*|(\\.|/)(test|spec))\\.(js|ts)$",
  // moduleFileExtensions: ["vue", "js"],
  // moduleNameMapper: {
  //   "^@/(.*)$": "<rootDir>/src/$1",
  // },
  coveragePathIgnorePatterns: ["/node_modules/", "/tests/"],
  coverageReporters: ["text", "json-summary"],
  testEnvironmentOptions: {
    customExportConditions: ["node", "node-addons"],
  },
};
