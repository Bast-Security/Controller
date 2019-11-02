defmodule Bast.Middleman.REST do
  @moduledoc """
  This is the HTTP server responsible for handling HTTP
  REST messages sent to the middleman.
  """

  import Plug.Conn

  def init(options) do
    options
  end

  def call(conn, _options) do
    conn
    |> put_resp_content_type("text/plain")
    |> send_resp(200, "Hello World!\n")
  end
end
