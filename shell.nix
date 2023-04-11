{ pkgs ? import <nixpkgs> {} }:

let
	present = pkgs.buildGoModule {
		name = "present";
		src = pkgs.fetchFromGitHub {
			owner = "golang";
			repo = "tools";
			rev = "v0.8.0";
			sha256 = "1lfz2wwd7jjzr3b3q7z91as22jic9rmqfw9wicj6bhr19swr5jvb";
		};
		subPackages = [ "cmd/present" ];
		vendorSha256 = "1289wjw0nxwwycmglvmds9djphy8gkamccb14sqzmf698dm83sdj";
	};
in

pkgs.mkShell {
	buildInputs = with pkgs; [
		present
		go
		deno
	];
}
