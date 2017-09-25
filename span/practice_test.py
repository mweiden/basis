import unittest

import practice


class TestStockSpan(unittest.TestCase):

    def setUp(self):
        pass

    def test_stock_spam(self):
        self.assertEquals(
            practice.calculate_span([1, 2, 3, 2, 6]),
            [1, 2, 3, 1, 5]
        )