defmodule BastTest do
  use ExUnit.Case
  doctest Bast

  test "greets the world" do
    assert Bast.hello() == :world
  end
end
