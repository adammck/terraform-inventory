class TerraformInventory < Formula
  homepage "https://github.com/adammck/terraform-inventory"
  head "https://github.com/adammck/terraform-inventory.git"

  # Update these when a new version is released
  url "https://github.com/adammck/terraform-inventory/archive/v0.6.1.tar.gz"
  sha256 "9cdcbc5ce4247b72bb72923d01246f51252a88908d760d766daeac51dd8feed9"

  depends_on "go" => :build

  def install
    ENV["GOPATH"] = buildpath

    # Move the contents of the repo (which are currently in the buildpath) into
    # a go-style subdir, so we can build it without spewing deps everywhere.
    app_path = buildpath/"src/github.com/adammck/terraform-inventory"
    app_path.install Dir["*"]

    # Fetch the deps (into our temporary gopath) and build
    cd "src/github.com/adammck/terraform-inventory" do
      system "go", "get"
      system "go", "build", "-ldflags", "-X main.build_version '#{version}'"
    end

    # Install the resulting binary
    bin.install "bin/terraform-inventory"
  end

  test do
    system "#{bin}/terraform-inventory", "version"
  end
end
