# Apache SkyWalking release guide

This documentation guides the release manager to release the SkyWalking CLI in the Apache Way, and also helps people to
check the release for vote.

## Add your GPG public key

1. Add your GPG public key into [SkyWalking GPG KEYS](https://dist.apache.org/repos/dist/release/skywalking/KEYS) file,
you can do this only if you are a committer, use your Apache id and password to log into the svn, and update the file.
**DO NOT override the existed `KEYS` file.**

1. Upload your GPG public key to public GPG site, such as [MIT's site](http://pgp.mit.edu:11371/). This site should be in 
Apache maven staging repository check list.

## Build and sign the source code package

```shell
export VERSION=<the version to release>
git tag v${VERSION}
make release
```

The `skywalking-cli-${VERSION}-bin.tgz`, `skywalking-cli-${VERSION}-src.tgz`, and their corresponding `asc`, `sha512`
should be generated in the directory automatically. 

## Upload to Apache svn

1. Check out the [dist dev repo](https://dist.apache.org/repos/dist/dev/skywalking), e.g. `svn checkout https://dist.apache.org/repos/dist/dev/skywalking`
1. Create a folder named with the release version and round, prefixed by `cli/`, e.g. `cd skywalking && mkdir -p cli/$VERSION`
1. Copy all the packages to the folder, as well as their `.asc` and `.sha512`
1. Commit the changes to svn

## Make the internal announcement

Send an announcement email to dev@ mailing list.

```text
Subject: [ANNOUNCEMENT] SkyWalking CLi $VERSION test build available

Content:

The test build of SkyWalking CLI $VERSION is now available.

This is our Apache release.
We welcome any comments you may have, and will take all feedback into
account if a quality vote is called for this build.

Release notes:

 * https://github.com/apache/skywalking-cli/blob/v$VERSION/CHANGES.md

Release Candidate:

 * https://dist.apache.org/repos/dist/dev/skywalking/cli/$VERSION
 * sha512 checksums
   - sha512xxxxyyyzzz apache-skywalking-cli-bin-x.x.x.tgz
   - sha512xxxxyyyzzz apache-skywalking-cli-src-x.x.x.tgz

Release Tag :

 * (Git Tag) v$VERSION

Release Commit Hash :

 * https://github.com/apache/skywalking-cli/tree/<Git Commit Hash>

Keys to verify the Release Candidate :

 * http://pgp.mit.edu:11371/pks/lookup?op=get&search=0x8BD99F552D9F33D7 corresponding to kezhenxu94@apache.org

Guide to build the release from source :

 * https://github.com/apache/skywalking-cli/blob/v$VERSION/docs/How-to-release.md

A vote regarding the quality of this test build will be initiated
within the next couple of days.
```

## Wait at least 48 hours for test responses

Any PMC, committer or contributor can test features for releasing, and feedback.
Based on that, PMC will decide whether to start a vote or not.

## Call for vote in dev@ mailing list

Call for vote in `dev@skywalking.apache.org`

```text
Subject: [VOTE] Release Apache SkyWalking CLI version $VERSION

Content:

Hi All,
This is a call for vote to release Apache SkyWalking CLI version $VERSION.

Release notes:

 * https://github.com/apache/skywalking-cli/blob/v$VERSION/CHANGES.md

Release Candidate:

 * https://dist.apache.org/repos/dist/dev/skywalking/cli/$VERSION
 * sha512 checksums
   - sha512xxxxyyyzzz apache-skywalking-cli-src-x.x.x.tgz
   - sha512xxxxyyyzzz apache-skywalking-cli-bin-x.x.x.tgz

Release Tag :

 * (Git Tag) v$VERSION

Release Commit Hash :

 * https://github.com/apache/skywalking-cli/tree/<Git Commit Hash>

Keys to verify the Release Candidate :

 * https://dist.apache.org/repos/dist/release/skywalking/KEYS

Guide to build the release from source :

 * https://github.com/apache/skywalking-cli/blob/v$VERSION/docs/How-to-release.md

Voting will start now (xxxx date) and will remain open for at least 72 hours, all PMC members are required to give their votes.
[ ] +1 Release this package.
[ ] +0 No opinion.
[ ] -1 Do not release this package because....
```

## Vote Check

All PMC members and committers should check these before voting +1.

1. Features test.
1. All artifacts in staging repository are published with `.asc`, `.md5`, and `sha` files
1. Source codes and distribution packages (`apache-skywalking-cli-{src,bin}-$VERSION.tgz`)
are in `https://dist.apache.org/repos/dist/dev/skywalking/cli/$VERSION` with `.asc`, `.sha512`.
1. `LICENSE` and `NOTICE` are in source codes and distribution package.
1. Check `shasum -c apache-skywalking-cli-{src,bin}-$VERSION.tgz.sha512`
1. Build distribution from source code package by following this [the build guide](#build-and-sign-the-source-code-package).
1. Licenses check, `make license`;

Vote result should follow these.

1. PMC vote is +1 binding, all others is +1 no binding.
1. Within 72 hours, you get at least 3 (+1 binding), and have more +1 than -1. Vote pass. 


## Publish release

1. Move source codes tar balls and distributions to `https://dist.apache.org/repos/dist/release/skywalking/`.

```shell
export SVN_EDITOR=vim
svn mv https://dist.apache.org/repos/dist/dev/skywalking/cli/$VERSION https://dist.apache.org/repos/dist/release/skywalking/cli
# ....
# enter your apache password
# ....

```

1. Public download links of source and distribution tar/zip locate in `http://www.apache.org/dyn/closer.cgi/skywalking/cli/$VERSION/`.
We only publish Apache mirror path as release info.

1. Public asc and sha512 locate in `https://www.apache.org/dist/skywalking/cli/$VERSION/xxx`

1. Public KEYS pointing to  `https://www.apache.org/dist/skywalking/KEYS`

1. Send ANNOUNCEMENT email to `dev@skywalking.apache.org`.

```text
Subject: [ANNOUNCEMENT] Release Apache SkyWalking CLI version $VERSION

Content:

Hi all,

The Apache SkyWalking Team is glad to announce the release of Apache SkyWalking CLI $VERSION.

SkyWalking CLI: Apache SkyWalking CLI https://skywalking.apache.org/

SkyWalking: APM (application performance monitor) tool for distributed systems, 
especially designed for microservices, cloud native and container-based (Docker, Kubernetes, Mesos) architectures. 

Vote Thread: 

Download Links : http://skywalking.apache.org/downloads/

Release Notes : https://github.com/apache/skywalking-cli/blob/v$VERSION/CHANGES.md

Website: http://skywalking.apache.org/

SkyWalking Resources:
- Issue: https://github.com/apache/skywalking/issues
- Mailing list: dev@skywalkiing.apache.org
- Documents: https://github.com/apache/skywalking/blob/v$VERSION/docs/README.md


- Apache SkyWalking Team
```

1. Update website download page. http://skywalking.apache.org/downloads/ , including the new download source, distribution, sha512, asc and document
links. Links could be found by following above rules.

1. Add a release event on website homepage and even page. Announce the public release with changelog or key features.
