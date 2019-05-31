import unittest
#from unittest.mock import patch
import mock
from users import get_users
from mock import patch
class BasicTests(unittest.TestCase):
    @patch('users.requests.post')  # Mock 'requests' module 'get' method.
    def test_request_response_with_decorator(self, mock_get):
        """Mocking using a decorator"""
        mock_get.return_value.status_code = 200# Mock status code of response.
        response = get_users()
        # Assert that the request-response cycle completed successfully.
        self.assertEqual(response.status_code, 200)
if __name__ == "__main__":
    unittest.main()
