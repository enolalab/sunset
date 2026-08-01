import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'Sunset',
  tagline: 'Source code documentation, indexed by tree-sitter',
  favicon: 'img/logo.svg',
  url: 'https://sunset.enolalab.com',
  baseUrl: '/',
  organizationName: 'enolalabs',
  projectName: 'sunset',
  onBrokenLinks: 'throw',
  i18n: {defaultLocale: 'en', locales: ['en']},
  presets: [
    [
      'classic',
      {
        docs: {
          routeBasePath: '/',
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/enolalabs/sunset/edit/main/website/',
        },
        blog: false,
        theme: {customCss: './src/css/custom.css'},
      } satisfies Preset.Options,
    ],
  ],
  themeConfig: {
    image: 'img/logo.svg',
    navbar: {
      title: 'Sunset',
      logo: {alt: 'Sunset logo', src: 'img/logo.svg'},
      items: [
        {type: 'docSidebar', sidebarId: 'docs', position: 'left', label: 'Documentation'},
        {href: 'https://enolalab.com', label: 'Enolalab', position: 'right'},
        {href: 'https://github.com/enolalabs/sunset', label: 'GitHub', position: 'right'},
        {href: 'https://github.com/enolalabs/sunset/releases', label: 'Releases', position: 'right'},
        {href: 'https://github.com/enolalabs/sunset/packages', label: 'Catalogue', position: 'right'},
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {title: 'Docs', items: [{label: 'Getting started', to: '/getting-started/installation'}, {label: 'CLI reference', to: '/cli/overview'}]},
        {title: 'Project', items: [{label: 'GitHub', href: 'https://github.com/enolalabs/sunset'}, {label: 'Releases', href: 'https://github.com/enolalabs/sunset/releases'}, {label: 'Catalogue', href: 'https://github.com/enolalabs/sunset/packages'}]},
        {title: 'Enolalab', items: [{label: 'enolalab.com', href: 'https://enolalab.com'}]},
      ],
      copyright: `Copyright ${new Date().getFullYear()} Enolalab.`,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
