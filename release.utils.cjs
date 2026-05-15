const ALL_SCOPES = ['api', 'ui', 'cli', 'chart-api', 'chart-ui', 'repo'];

function buildReleaseRules(scope) {
  return [
    ...ALL_SCOPES.filter(s => s !== scope).map(s => ({ scope: s, release: false })),
    { type: 'feat', scope: scope, release: 'minor' },
    { type: 'fix', scope: scope, release: 'patch' },
    { type: 'perf', scope: scope, release: 'patch' },
    { type: 'revert', scope: scope, release: 'patch' },
    { type: 'build', scope: scope, release: 'patch' },
    { type: 'chore', scope: scope, release: 'patch' },
    { type: 'refactor', scope: scope, release: 'patch' },
    { type: 'test', scope: scope, release: 'patch' },
  ];
}

function buildReleaseNotesConfig(scope) {
  return {
    preset: 'conventionalcommits',
    presetConfig: {
      types: [
        { type: 'feat',     section: 'Features' },
        { type: 'fix',      section: 'Bug Fixes' },
        { type: 'perf',     section: 'Performance' },
        { type: 'revert',   section: 'Reverts' },
        { type: 'build',    section: 'Build' },
        { type: 'chore',    section: 'Chores' },
        { type: 'refactor', section: 'Refactors' },
        { type: 'test',     section: 'Tests' },
      ],
    },
    writerOpts: {
      transform: (commit) => {
        if (commit.scope !== scope) return false;
        return commit;
      },
    },
  };
}

module.exports = { buildReleaseRules, buildReleaseNotesConfig };
