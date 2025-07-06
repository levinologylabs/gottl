import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  base: process.env.VITEPRESS_BASE || '/',
  title: "Gottl",
  description: "A Go JSON API Starter Kit",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Quick Start', link: '/about/quick-start' }
    ],

    sidebar: [
      {
        text: 'About',
        items: [
          { text: 'What is Gottl', link: '/about/what-is-gottl' },
          { text: 'Quick Start', link: '/about/quick-start' }
        ]
      },
      {
        text: 'User Guide',
        items: [
          { text: 'Project Layout', link: '/user-guide/project-layout' },
          { text: 'Configuration', link: '/user-guide/configuration' },
          { text: 'Database and Migrations', link: '/user-guide/database-and-migrations' },
          { text: 'Open Telemetry', link: '/user-guide/open-telemetry' },
          { text: 'Scaffolding APIs', link: '/user-guide/scaffolding' },
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/levinologylabs/gottl' }
    ]
  }
})
