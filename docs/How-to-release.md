# Apache SkyWalking CLI Release Guide

This documentation guides the release manager to release the SkyWalking CLI in the Apache Way, and also helps people to
check the release for vote.

## Prerequisites

1. Close(if finished, or move to next milestone otherwise) all issues in the current milestone
   from [skywalking-cli](https://github.com/apache/skywalking-cli/milestones)
   and [skywalking](https://github.com/apache/skywalking/milestones), create a new milestone if needed.
2. Update [CHANGES.md](../CHANGES.md).

## Add your GPG public key to Apache svn

1. Upload your GPG public key to a public GPG site, such as [MIT's site](http://pgp.mit.edu:11371/).

1. Log in [id.apache.org](https://id.apache.org/) and submit your key fingerprint.

1. Add your GPG public key into [SkyWalking GPG KEYS](https://dist.apache.org/repos/dist/release/skywalking/KEYS)
   file, **you can do this only if you are a PMC member**. You can ask a PMC member for help. **DO NOT override the
   existed `KEYS` file content, only append your key at the end of the file.**

## Build and sign the source code package

```shell
export VERSION=<the version to release>
git clone --recurse-submodules git@github.com:apache/skywalking-cli && cd skywalking-cli
git tag -a "$VERSION" -m "Release Apache SkyWalking-CLI $VERSION"
git push --tags
make clean && make release
```

**In total, six files should be automatically generated in the directory**: `skywalking-cli-${VERSION}-bin.tgz`
, `skywalking-cli-${VERSION}-src.tgz`, and their corresponding `asc`, `sha512` files.

## Upload to Apache svn

```bash
svn co https://dist.apache.org/repos/dist/dev/skywalking/
mkdir -p skywalking/cli/"$VERSION"
cp skywalking-cli/skywalking*.tgz skywalking/cli/"$VERSION"
cp skywalking-cli/skywalking*.tgz.asc skywalking/cli/"$VERSION"
cp skywalking-cli/skywalking-cli*.tgz.sha512 skywalking/cli/"$VERSION"

cd skywalking/cli && svn add "$VERSION" && svn commit -m "Draft Apache SkyWalking-CLI release $VERSION"
```

## Make the internal announcement

Send an announcement email to dev@ mailing list, **please check all links before sending the email**.

```text
Subject: [ANNOUNCEMENT] SkyWalking CLi $VERSION test build available

Content:

The test build of SkyWalking CLI $VERSION is now available.

We welcome any comments you may have, and will take all feedback into
account if a quality vote is called for this build.

Release notes:

 * https://github.com/apache/skywalking-cli/blob/$VERSION/CHANGES.md

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

 * https://github.com/apache/skywalking-cli/blob/$VERSION/docs/How-to-release.md

A vote regarding the quality of this test build will be initiated
within the next couple of days.
```

## Wait at least 48 hours for test responses

Any PMC, committer or contributor can test features for releasing, and feedback. Based on that, PMC will decide whether
to start a vote or not.

## Call for vote in dev@ mailing list

Call for vote in `dev@skywalking.apache.org`, **please check all links before sending the email**.

```text
Subject: [VOTE] Release Apache SkyWalking CLI version $VERSION

Content:

Hi the SkyWalking Community:
This is a call for vote to release Apache SkyWalking CLI version $VERSION.

Release notes:

 * https://github.com/apache/skywalking-cli/blob/$VERSION/CHANGES.md

Release Candidate:

 * https://dist.apache.org/repos/dist/dev/skywalking/cli/$VERSION
 * sha512 checksums
   - sha512xxxxyyyzzz skywalking-cli-x.x.x-src.tgz
   - sha512xxxxyyyzzz skywalking-cli-x.x.x-bin.tgz

Release Tag :

 * (Git Tag) v$VERSION

Release Commit Hash :

 * https://github.com/apache/skywalking-cli/tree/<Git Commit Hash>

Keys to verify the Release Candidate :

 * https://dist.apache.org/repos/dist/release/skywalking/KEYS

Guide to build the release from source :

 * https://github.com/apache/skywalking-cli/blob/$VERSION/docs/How-to-release.md

Voting will start now and will remain open for at least 72 hours, all PMC members are required to give their votes.

[ ] +1 Release this package.
[ ] +0 No opinion.
[ ] -1 Do not release this package because....

Thanks.

[1] https://github.com/apache/skywalking-cli/blob/master/docs/How-to-release.md#vote-check
```

## Vote Check

All PMC members and committers should check these before voting +1:

1. Features test.
1. All artifacts in staging repository are published with `.asc`, `.md5`, and `sha` files.
1. Source codes and distribution packages (`skywalking-cli-$VERSION-{src,bin}.tgz`)
   are in `https://dist.apache.org/repos/dist/dev/skywalking/cli/$VERSION` with `.asc`, `.sha512`.
1. `LICENSE` and `NOTICE` are in source codes and distribution package.
1. Check `shasum -c skywalking-cli-$VERSION-{src,bin}.tgz.sha512`.
1. Check `gpg --verify skywalking-cli-$VERSION-{src,bin}.tgz.asc skywalking-cli-$VERSION-{src,bin}.tgz`.
1. Build distribution from source code package by following
   this [the build guide](#build-and-sign-the-source-code-package).
1. Licenses check, `make license`.

Vote result should follow these:

1. PMC vote is +1 binding, all others is +1 no binding.

1. Within 72 hours, you get at least 3 (+1 binding), and have more +1 than -1. Vote pass.

1. **Send the closing vote mail to announce the result**. When count the binding and no binding votes, please list the
   names of voters. An example like this:

   ```
   [RESULT][VOTE] Release Apache SkyWalking CLI version $VERSION
   
   3 days passed, we’ve got ($NUMBER) +1 bindings (and ... +1 non-bindings):
   
   (list names)
   +1 bindings:
   xxx
   ...
      
   +1 non-bindings:
   xxx
   ...
    
   Thank you for voting, I’ll continue the release process.
   ```

## Publish release

1. Move source codes tar balls and distributions to `https://dist.apache.org/repos/dist/release/skywalking/`, **you can
   do this only if you are a PMC member**.

    ```shell
    svn mv https://dist.apache.org/repos/dist/dev/skywalking/cli/$VERSION https://dist.apache.org/repos/dist/release/skywalking/cli -m"Release SkyWalking CLI $VERSION"
    ```

2. Push Docker images.

In order to publish the Docker images, you
have [to enable BuildKit in the Docker daemon](https://docs.docker.com/develop/develop-images/build_enhancements/.

```shell
make docker.push
```

1. Refer to the previous [PR](https://github.com/apache/skywalking-website/pull/212), update the event and download
   links on the website.

1. Update [Github release page](https://github.com/apache/skywalking-cli/releases), follow the previous convention.

1. Send ANNOUNCE email to `dev@skywalking.apache.org` and `announce@apache.org`, the sender should use his/her Apache
   email account, **please check all links before sending the email**.

    ```
    Subject: [ANNOUNCEMENT] Apache SkyWalking CLI $VERSION Released

    Content:

    Hi the SkyWalking Community

    On behalf of the SkyWalking Team, I’m glad to announce that SkyWalking CLI $VERSION is now released.

    SkyWalking CLI: A Command Line Interface for Apache SkyWalking.

    SkyWalking: APM (application performance monitor) tool for distributed systems, especially designed for microservices, cloud native and container-based (Docker, Kubernetes, Mesos) architectures.

    Download Links: http://skywalking.apache.org/downloads/

    Release Notes : https://github.com/apache/skywalking-cli/blob/$VERSION/CHANGES.md

    Website: http://skywalking.apache.org/

    SkyWalking CLI Resources:
    - Issue: https://github.com/apache/skywalking/issues
    - Mailing list: dev@skywalking.apache.org
    - Documents: https://github.com/apache/skywalking-cli/blob/$VERSION/README.md
    
    The Apache SkyWalking Team
    ```

1. Add a release on [the report system](https://reporter.apache.org/addrelease.html?skywalking) (**you can do this only
   if you are a PMC member**).

## Remove Unnecessary Releases

Please remember to remove all unnecessary releases in the mirror
svn (https://dist.apache.org/repos/dist/release/skywalking/), if you don't recommend users to choose those version. For
example, you have removed the download and documentation links from the website. If they want old ones, the Archive
repository has all of them.
