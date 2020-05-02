defmodule Bast.MixProject do
  use Mix.Project

  def project do
    [
      app: :bast,
      version: "0.1.0",
      elixir: "~> 1.9",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger],
      mod: {Bast.Application, []}
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      plug_cowboy: "~> 2.0",
      poison: "~> 3.1",
      ecto_sql: "~> 3.2",
      postgrex: "~> 0.15"
    ]
  end
end
