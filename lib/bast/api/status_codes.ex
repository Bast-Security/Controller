defmodule Bast.Api.Status do
  def ok, do: { 200, "Ok\n" }

  def created, do: { 201, "Created\n" }

  def malformed_request, do: { 422, "Malformed Request.\n" }

  def unauthorized, do: { 401, "Unauthorized\n" }

  def forbidden, do: { 403, "Forbidden\n" }

  def not_found, do: { 404, "Not Found\n" }

  def not_implemented, do: { 501, "Not Implemented\n" }

  def status_from_crud({:ok, _struct}, :create), do: created()

  def status_from_crud({:ok, _struct} ), do: ok()

  def status_from_crud({:error, _changeset}), do: malformed_request()

  def send_status({status, body}, conn) do
    Plug.Conn.send_resp(conn, status, body)
  end
end
