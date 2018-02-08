import tempfile
import unittest

import yaml

from py import release


class ReleaseTest(unittest.TestCase):
  def test_update_values(self):
    with tempfile.NamedTemporaryFile(delete=False) as hf:
      hf.write("""# Test file
image: gcr.io/image:latest

## Install Default RBAC roles and bindings
rbac:
  install: false
  apiVersion: v1beta1""")
      values_file = hf.name

    release.update_values(hf.name, "gcr.io/image:v20171019")

    with open(values_file) as hf:
      output = hf.read()

      expected = """# Test file
image: gcr.io/image:v20171019

## Install Default RBAC roles and bindings
rbac:
  install: false
  apiVersion: v1beta1"""
      self.assertEquals(expected, output)

  def test_update_chart_file(self):
    with tempfile.NamedTemporaryFile(delete=False) as hf:
      hf.write("""
name: tf-job-operator-chart
home: https://github.com/jlewi/mlkube.io
version: 0.1.0
appVersion: 0.1.0
""")
      chart_file = hf.name

    release.update_chart(chart_file, "v20171019")

    with open(chart_file) as hf:
      output = yaml.load(hf)
    expected = {
        "name": "tf-job-operator-chart",
        "home": "https://github.com/jlewi/mlkube.io",
        "version": "0.1.0-v20171019",
        "appVersion": "0.1.0-v20171019",
    }
    self.assertEquals(expected, output)


if __name__ == "__main__":
  unittest.main()
