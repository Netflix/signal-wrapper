#!/usr/bin/env python3

import os
import unittest

start_test=int(os.environ['start_test'])

start_trap=int(os.environ['start_trap'])
start_shutdown=int(os.environ['start_shutdown'])

finish_shutdown=int(os.environ['finish_shutdown'])
finish_trap=int(os.environ['finish_trap'])


class TestSignalWatcher(unittest.TestCase):
	def test_p1(self):
		self.assertTrue(start_trap - start_test <= 1)
	def test_p2(self):
		# This should be about 5 seconds, based on the sleep before kill in test.sh
		self.assertTrue(start_shutdown - start_trap <= 5)
	def test_p3(self):
		self.assertTrue(finish_shutdown - start_shutdown <= 9)
	def test_p4(self):
		self.assertTrue(finish_shutdown - finish_trap <= 1)

if __name__ == '__main__':
	unittest.main()
