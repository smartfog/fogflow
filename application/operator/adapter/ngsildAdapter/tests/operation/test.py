import unittest
from main import get_users


class BasicTests(unittest.TestCase):
    def test_request_response(self):
        response = get_users()
        # Assert that the request-response cycle completed successfully with status code 200.
        self.assertEqual(response.status_code, 200)
if __name__ == "__main__":
    unittest.main()
