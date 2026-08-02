import type {SidebarItem} from '@docusaurus/plugin-content-docs';

const sidebars: Record<string, SidebarItem[]> = {
  docs: [
    'intro',
    {type: 'category', label: 'Getting started', items: ['getting-started/installation', 'getting-started/quick-start']},
    {type: 'category', label: 'CLI', items: ['cli/overview', 'cli/parse', 'cli/update', 'cli/languages', 'cli/version', 'cli/clean']},
    {type: 'category', label: 'Concepts', items: ['concepts/supported-languages', 'concepts/scanning-and-exclusions', 'concepts/output-format', 'concepts/caching']},
    {type: 'category', label: 'Go API', items: ['go-api/overview', 'go-api/parsing', 'go-api/tree-traversal', 'go-api/language-detection']},
    {type: 'category', label: 'Maintainers', items: ['maintainers/contributing', 'maintainers/releases']},
    {type: 'category', label: 'Reference', items: ['reference/limitations']},
  ],
};

export default sidebars;
