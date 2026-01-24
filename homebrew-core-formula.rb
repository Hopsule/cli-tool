# This is the Homebrew Core formula for hopsule
# To be submitted as a PR to https://github.com/Homebrew/homebrew-core

class Hopsule < Formula
  desc "Decision & Memory Layer for AI teams - CLI tool"
  homepage "https://github.com/Hopsule/cli-tool"
  url "https://github.com/Hopsule/cli-tool/archive/v0.4.4.tar.gz"
  sha256 "CHECKSUM_TO_BE_CALCULATED"
  license "MIT"
  head "https://github.com/Hopsule/cli-tool.git", branch: "main"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X main.version=#{version}
      -X main.commit=#{tap.user}
      -X main.date=#{time.iso8601}
    ]

    system "go", "build", *std_go_args(ldflags: ldflags, output: bin/"hopsule"), "./cmd/decision"

    # Install shell completions (if added later)
    # generate_completions_from_executable(bin/"hopsule", "completion")
  end

  test do
    assert_match "decision version", shell_output("#{bin}/hopsule --version")
    
    # Test config creation
    system bin/"hopsule", "config"
    assert_predicate testpath/".decision-cli/config.yaml", :exist?
  end
end
