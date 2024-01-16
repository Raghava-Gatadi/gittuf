## gittuf policy add-rule

Add a new rule to a policy file

### Synopsis

This command allows users to add a new rule to the specified policy file. By default, the main policy file is selected. Note that authorized keys can be specified from disk, from the GPG keyring using the "gpg:<fingerprint>" format, or as a Sigstore identity as "fulcio:<identity>::<issuer>".

```
gittuf policy add-rule [flags]
```

### Options

```
      --authorize-key stringArray   authorized public key for rule
  -h, --help                        help for add-rule
      --policy-name string          policy file to add rule to (default "targets")
      --rule-name string            name of rule
      --rule-pattern stringArray    patterns used to identify namespaces rule applies to
```

### Options inherited from parent commands

```
  -k, --signing-key string   signing key to use to sign policy file
      --verbose              enable verbose logging
```

### SEE ALSO

* [gittuf policy](gittuf_policy.md)	 - Tools to manage gittuf policies
