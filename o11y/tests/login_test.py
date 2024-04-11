from observability import login

def test_login_successful():
    login_string = "3cfa8b505c2a2a2e2b54bb6081c8d9fcefd5b836@example.grafana.com"
    result = login(login_string)
    assert result == True

def test_login_failed():
    token = "invalid_token"
    result = login(token)
    assert result == False