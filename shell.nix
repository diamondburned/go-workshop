let sources = import ./nix/sources.nix;
in

{ pkgs ? import sources.nixpkgs {} }:

let
	present = pkgs.buildGoModule {
		name = "present";
		src = pkgs.runCommand "present-src" {
			src = pkgs.fetchFromGitHub {
				owner = "golang";
				repo = "tools";
				rev = "v0.8.0";
				sha256 = "1lfz2wwd7jjzr3b3q7z91as22jic9rmqfw9wicj6bhr19swr5jvb";
			};
		} ''
			mkdir -p $out
			cp -r $src/* $out
			chmod -R +w $out
			cp -f ${./static}/* $out/cmd/present/static
			cp -f ${./templates}/* $out/cmd/present/templates
		'';
		vendorSha256 = "1289wjw0nxwwycmglvmds9djphy8gkamccb14sqzmf698dm83sdj";
		subPackages = [ "cmd/present" ];
	};

	httpie = pkgs.buildGoModule {
		name = "httpie";
		src = pkgs.fetchFromGitHub {
			owner = "nojima";
			repo = "httpie-go";
			rev = "v0.7.0";
			sha256 = "0p67wmfvdfyl28i7sbwmnhr869077l48a1v9icmn05nkim59gs9i";
		};
		vendorSha256 = "1ib58bs47jjyj5wykcn4f3ks0jmffpb7qnk74d99qx6pli9p6y68";
		postInstall = ''mv $out/bin/ht $out/bin/httpie'';
	};
in

pkgs.mkShell {
	buildInputs = with pkgs; [
		present
		httpie
		niv
		go
		deno
		nodePackages.prettier
		(pkgs.writeShellScriptBin "deno-run" "deno run -A $@")
	];
}
