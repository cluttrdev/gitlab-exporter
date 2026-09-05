#!/bin/bash
set -e

gitlab-rails runner - <<'EOS'
user = User.find_by_username('root')
token = user.personal_access_tokens.find_by(name: 'setup-token')
unless token
  token = user.personal_access_tokens.new(
    name: 'setup-token',
    scopes: [:api],
    expires_at: 365.days.from_now
  )
  token.set_token('glpat-setup-root-token')
  token.save!
end
EOS

echo "Post Reconfigure Script successfully executed"
