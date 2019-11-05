defmodule Bast.Controller.MixProject do
  use Mix.Project

  def project do
    [
      app: :controller,
      version: "0.1.0",
      elixir: "~> 1.7",
      start_permanent: Mix.env() == :prod,
      deps: deps()
    ]
  end

  # Run "mix help compile.app" to learn about applications.
  def application do
    [
      extra_applications: [:logger],
      mod: {Bast.Controller.Application, []}
    ]
  end

  # Run "mix help deps" to learn about dependencies.
  defp deps do
    [
      plug_cowboy: "~> 2.0",
      ecto_sql: "~> 3.2",
      myxql: "~> 0.2.10",
      jason: "~> 1.1"
    ]
  end
end
