{ lib
, buildGoModule
}:

buildGoModule {
  pname = "uci-ml-repository-dataset-metadata";
  version = "unstable-2023-09-15";

  src = ./.;

  vendorHash = "sha256-htG3NNDj080M6vY2+xm81xoCC0QT7w714bY5CgJ/SzI=";

  ldflags = [ "-s" "-w" ];

  meta = with lib; {
    description = "Download and parse metadata on all datasets in the new UCI ML repository";
    homepage = "https://github.com/Leixb/UCI-ML-Repository-dataset-metadata/";
    # license = licenses.unfree; # FIXME: nix-init did not found a license
    maintainers = with maintainers; [ ];
    mainProgram = "uciml";
  };
}
