{
  dockerTools,
  runCommandLocal,
  upside-matrix,
  benchmark-testcase,
}:
let
  patched-upsided = upside-matrix.upsided.overrideAttrs (oldAttrs: {
    patches = oldAttrs.patches or [ ] ++ [ ./testground-upsided.patch ];
  });
in
let
  tmpDir = runCommandLocal "tmp" { } ''
    mkdir -p $out/tmp/
  '';
in
dockerTools.buildLayeredImage {
  name = "upside-testground";
  created = "now";
  contents = [
    benchmark-testcase
    patched-upsided
    tmpDir
  ];
  config = {
    Expose = [
      9090
      26657
      26656
      1317
      26658
      26660
      26659
      30000
    ];
    Cmd = [ "/bin/stateless-testcase" ];
    Env = [ "PYTHONUNBUFFERED=1" ];
  };
}
