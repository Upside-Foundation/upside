{
  pkgs,
  config,
  upside ? (import ../. { inherit pkgs; }),
}:
rec {
  start-upside = pkgs.writeShellScriptBin "start-upside" ''
    # rely on environment to provide upsided
    export PATH=${pkgs.test-env}/bin:$PATH
    ${../scripts/start-upside} ${config.upside-config} ${config.dotenv} $@
  '';
  start-geth = pkgs.writeShellScriptBin "start-geth" ''
    export PATH=${pkgs.test-env}/bin:${pkgs.go-ethereum}/bin:$PATH
    source ${config.dotenv}
    ${../scripts/start-geth} ${config.geth-genesis} $@
  '';
  start-scripts = pkgs.symlinkJoin {
    name = "start-scripts";
    paths = [
      start-upside
      start-geth
    ];
  };
}
