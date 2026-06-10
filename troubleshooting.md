import unittest
from social_assembly_docs import DeploymentGuide

class TestIntegration(unittest.TestCase):
    def test_deployment_guide_generation(self):
        guide = DeploymentGuide()
        self.assertTrue(guide.generate_guide())

    def test_nginx_config_parsing(self):
        config_file = '/opt/axentx/social-assembly/docs/nginx-example.conf'
        with open(config_file, 'r') as f:
            contents = f.read()
            self.assertTrue(DeploymentGuide.parse_nginx_config(contents))

    def test_troubleshooting_section(self):
        guide = DeploymentGuide()
        self.assertTrue(guide.get_troubleshooting_section())

if __name__ == '__main__':
    unittest.main()