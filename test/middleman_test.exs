defmodule MiddlemanTest do
  use ExUnit.Case
  doctest Middleman

  test "greets the world" do
    assert Middleman.hello() == :world
  end
end
