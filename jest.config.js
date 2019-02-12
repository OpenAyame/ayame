module.exports = {
  preset: 'ts-jest',
  moduleFileExtensions: [
    'ts',
    'js'
  ],
  transform: {
    '^.+\\.tsx?$': 'ts-jest'
  },
  testMatch: [
    '**/test/**/*.test.(ts|js)'
  ],
  testEnvironment: 'node'
};
