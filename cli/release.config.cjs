const { buildReleaseRules } = require('../release.utils.cjs');

module.exports = {
  branches: ['main'],
  tagFormat: 'cli-v${version}',
  plugins: [
    ['@semantic-release/commit-analyzer', {
      preset: 'conventionalcommits',
      releaseRules: buildReleaseRules('cli'),
    }],
    ['@semantic-release/release-notes-generator', {
      preset: 'conventionalcommits',
    }],
    ['@semantic-release/github', {
      successComment: false,
      labels: false,
    }],
    ['@semantic-release/exec', {
      publishCmd: 'echo "${nextRelease.version}" > nextversion && echo "true" > released',
    }],
  ],
};
