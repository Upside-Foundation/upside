{
  pkgs ? import ../../nix { },
}:
let
  upsided = (pkgs.callPackage ../../. { });
in
upsided.overrideAttrs (oldAttrs: {
  patches = oldAttrs.patches or [ ] ++ [ ./broken-upsided.patch ];
})
