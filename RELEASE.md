Soon releasing will be as easy as
    git commit *version bump* ; git tag v0.0.20; git push

But until then,
    git commit *version bump* ; git tag v0.0.20; git push; git status #*to check for cleanliness*
    goreleaser release --clean

