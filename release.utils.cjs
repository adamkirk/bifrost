const ALL_SCOPES = ['api', 'ui', 'cli', 'chart-api', 'chart-ui', 'repo'];

function buildReleaseRules(scope) {
  return [
    ...ALL_SCOPES.filter(s => s !== scope).map(s => ({ scope: s, release: false })),
    { type: 'feat', scope: scope, release: 'minor' },
    { type: 'fix', scope: scope, release: 'patch' },
    { type: 'perf', scope: scope, release: 'patch' },
    { type: 'revert', scope: scope, release: 'patch' },
  ];
}

module.exports = { buildReleaseRules };
