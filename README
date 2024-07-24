# Filecrypt

Filecrypt is a file encryptor that packs files into a gzipped tarball in-memory,
and this is encrypted using NaCl via a scrypt-derived key.

It is influenced by the example in the "Practical Cryptography with Go."

## Usage

```text
filecrypt [-h] [-o filename] [-q] [-t] [-u] [-v] [-x] files...

        -h              Print this help message.

        -o filename     The filename to output. If an archive is being built,
                        this is the filename of the archive. If an archive is
                        being unpacked, this is the directory to unpack in.
                        If the tarball is being extracted, this is the path
                        to write the tarball.

                        Defaults:
                                   Pack: files.enc
                                 Unpack: .
                                Extract: files.tgz

        -q              Quiet mode. Only print errors and password prompt.
                        This will override the verbose flag.

	-t		List files in the archive. This acts like the list
			flag in tar.

        -u              Unpack the archive listed on the command line. Only
                        one archive may be unpacked.

        -v              Verbose mode. This acts like the verbose flag in
                        tar.

        -x              Extract a tarball. This will decrypt the archive, but
                        not decompress or unpack it.

Examples:
        filecrypt -o ssh.enc ~/.ssh
                Encrypt the user's OpenSSH directory to ssh.enc.

        filecrypt -o backup/ -u ssh.enc
                Restore the user's OpenSSH directory to the backup/
                directory.

        filecrypt -u ssh.enc
                Restore the user's OpenSSH directory to the current directory.
```
