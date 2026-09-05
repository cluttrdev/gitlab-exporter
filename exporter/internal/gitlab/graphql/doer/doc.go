// Package doer provides a thin, transport-level layer over genqlient-generated
// GraphQL clients that rewrites outgoing queries to strip version-gated
// selections based on the connected GitLab server's version.
package doer
