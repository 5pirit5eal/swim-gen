# Swim-RAG Frontend

Vue 3 + TypeScript frontend for the Swim Training Plan Generator application.

## Project Setup

```sh
npm install
```

### Development Commands

```sh
npm run dev          # Start development server with HMR
npm run build        # Type-check + build (runs in parallel)
npm run type-check   # TypeScript validation with vue-tsc
npm run test:unit    # Run Vitest tests in watch mode
npm run lint         # ESLint with auto-fix
npm run format       # Prettier formatting
npm run preview      # Preview production build locally
```

## V1 Development Todo List

### 🏗️ Project Setup & Architecture

- [x] ~~Initialize Vue 3 + TypeScript project with Vite~~
- [x] ~~Configure ESLint, Prettier, and VS Code settings~~
- [x] ~~Set up Pinia for state management~~
- [x] ~~Configure Vue Router for navigation~~
- [x] ~~Set up Vitest for unit testing~~
- [ ] **Plan application architecture and component structure**
  - [x] ~~Research Vue.js architecture patterns and best practices~~
  - [x] ~~Design component hierarchy and folder structure~~
  - [ ] **Create directory structure and placeholder files**
  - [ ] **Set up TypeScript type definitions**
  - [ ] **Configure API client structure**
  - [ ] **Set up Pinia stores structure**
- [ ] **Design API client structure for backend communication**
- [ ] **Set up environment variables for backend URL configuration**

### 🎨 UI/UX Design & Layout

- [ ] **Create main application layout component**
- [ ] **Design responsive header with navigation**
- [ ] **Implement footer with privacy/legal links**
- [ ] **Set up CSS design system (colors, typography, spacing)**
- [ ] **Create loading states and error handling UI patterns**

### 📝 Core Input Components

- [ ] **Build main text input form for training plan requests**
  - Free-form text area with proper validation
  - Character limits and input guidelines
  - Real-time input feedback
- [ ] **Create advanced settings panel**
  - Configuration options for training parameters
  - Collapsible/expandable design
  - Form validation and default values
- [ ] **Implement privacy settings controls**
  - Data donation opt-out checkbox
  - Clear privacy policy links
  - User consent management

### 🏊 Training Plan Features

- [ ] **Design training plan display component**
  - Structured display of generated plans
  - Readable formatting for exercises and sets
  - Clear organization by workout sections
- [ ] **Implement plan generation workflow**
  - Loading states during API calls
  - Error handling for failed requests
  - Success feedback and plan preview
- [ ] **Add plan customization options**
  - Edit generated plans before export
  - Save draft functionality (local storage)
  - Reset/regenerate options

### 📄 PDF Export System

- [ ] **Research PDF generation libraries for Vue.js**
  - Evaluate options (jsPDF, Puppeteer, html2pdf, etc.)
  - Consider server-side vs client-side generation
- [ ] **Implement PDF export functionality**
  - Format training plans for print
  - Include proper headers, footers, and branding
  - Handle different paper sizes and orientations
- [ ] **Add export options and preview**
  - PDF preview before download
  - Custom filename generation
  - Multiple export formats consideration for future

### 🔌 Backend Integration

- [ ] **Set up API client with proper TypeScript types**
  - HTTP client configuration (axios/fetch)
  - Request/response type definitions
  - Error handling and retry logic
- [ ] **Implement training plan generation API calls**
  - POST endpoint for plan requests
  - Proper request payload formatting
  - Response parsing and validation
- [ ] **Add API error handling and user feedback**
  - Network error recovery
  - Backend error message display
  - Offline state handling

### 🧪 Testing & Quality Assurance

- [ ] **Write unit tests for core components**
  - Input form validation tests
  - Plan display component tests
  - Store/state management tests
- [ ] **Implement integration tests**
  - End-to-end user workflows
  - API integration testing
  - PDF export functionality tests
- [ ] **Add accessibility testing and improvements**
  - Screen reader compatibility
  - Keyboard navigation support
  - Color contrast and visual accessibility

### 🚀 Deployment & Production

- [ ] **Configure production build optimization**
  - Bundle size analysis and optimization
  - Asset optimization (images, fonts)
  - Performance monitoring setup
- [ ] **Set up environment-specific configurations**
  - Development vs production API URLs
  - Feature flags for testing
  - Error reporting integration
- [ ] **Prepare for Cloud Run deployment**
  - Docker configuration
  - Static file serving
  - Health check endpoints

### 📚 Documentation & Learning

- [ ] **Document component APIs and usage patterns**
- [ ] **Create development workflow documentation**
- [ ] **Add inline code comments for learning**
- [ ] **Set up component storybook or documentation site**

### 🔮 Future Preparation (V2 Considerations)

- [ ] **Design authentication-ready architecture**
- [ ] **Plan for user session management**
- [ ] **Consider chat/history features in component design**
- [ ] **Prepare for multi-page application structure**

## Development Notes

### Learning Resources

- [Vue 3 Official Documentation](https://vuejs.org/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Vite Guide](https://vite.dev/guide/)
- [Pinia Documentation](https://pinia.vuejs.org/)
- [Vue Router Documentation](https://router.vuejs.org/)

### Architecture Decisions

- **Single Page Application**: V1 focuses on one main interface
- **Anonymous First**: No authentication complexity in initial version
- **Component-First**: Build reusable components for future versions
- **Type Safety**: Leverage TypeScript for better development experience
- **Modern Tooling**: Use latest Vue 3 Composition API and modern build tools

## Project Structure

```text
src/
├── components/          # Reusable UI components
├── views/              # Page-level components
├── stores/             # Pinia state management
├── router/             # Vue Router configuration
├── api/                # Backend API client (to be created)
├── types/              # TypeScript type definitions (to be created)
├── utils/              # Helper functions and utilities (to be created)
└── assets/             # Static assets (styles, images)
```
