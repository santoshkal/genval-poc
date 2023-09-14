          TAG=$(git describe --tags --abbrev=0)
          echo $TAG
          checksum=$(https://github.com/santoshkal/genval-poc/releases/download/$TAG/checksums.txt)
          echo $checksum
          cert=$(https://github.com/santoshkal/genval-poc/releases/download/$TAG/checksums.txt.pem)
          echo $cert
          sig=$(https://github.com/santoshkal/genval-poc/releases/download/$TAG/checksums.txt.sig)
          echo $sig


          cosign verify-blob \
          --certificate-identity-regexp='^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$' \
          --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
          --cert $cert --signature $sig $checksum
