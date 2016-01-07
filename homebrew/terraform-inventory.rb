class TerraformInventory < Formula
  homepage "https://github.com/adammck/terraform-inventory"
  head "https://github.com/adammck/terraform-inventory.git"

  # Update these when a new version is released
  url "https://github.com/adammck/terraform-inventory/archive/v0.6.tar.gz"
  sha1 "aa64cd2bd6f45d82f372f35ae6f7a5ec446a050a"

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
