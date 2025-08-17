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
- [x] ~~**Plan application architecture and component structure**~~
  - [x] ~~Research Vue.js architecture patterns and best practices~~
  - [x] ~~Design component hierarchy and folder structure~~
  - [x] ~~Create directory structure and placeholder files~~
  - [x] ~~Set up TypeScript type definitions with backend API alignment~~
  - [x] ~~Configure API client structure with error handling and timeouts~~
  - [x] ~~Set up Pinia stores structure~~
- [x] ~~**Design API client structure for backend communication**~~
- [x] ~~**Set up environment variables for backend URL configuration (via Vite proxy)**~~

### 🎨 UI/UX Design & Layout

- [x] ~~Create main application layout component~~
- [x] ~~Design responsive header with navigation~~
- [x] ~~Implement footer with privacy/legal links~~
- [x] ~~Set up CSS design system (colors, typography, spacing)~~
- [x] ~~Create loading states and error handling UI patterns~~
- [x] ~~~Implement tooltips as hoverabel question mark icon~~
- [x] ~~Add "i feel lucky" button right of the advanced settings button~~
- [x] Add icon and page name
- [ ] Translate page texts to german and implement multi-language in prep for v2**

### 📝 Core Input Components

- [x] ~~**Build main text input form for training plan requests**~~
  - [x] ~~Free-form text area with proper validation~~
  - [x] ~~Character limits and input guidelines~~
  - [x] ~~Real-time input feedback~~
- [x] ~~**Create advanced settings panel**~~
  - [x] ~~Configuration options for training parameters~~
  - [x] ~~Collapsible/expandable design~~
  - [x] ~~Form validation and default values~~
- [x] ~~**Implement privacy settings controls**~~
  - [x] ~~Data donation opt-out checkbox~~
  - [x] ~~Clear privacy policy links~~
  - [x] ~~User consent management~~

### 🏊 Training Plan Features

- [x] ~~**Design training plan display component**~~
  - [x] ~~Structured display of generated plans~~
  - [x] ~~Readable formatting for exercises and sets~~
  - [x] ~~Clear organization by workout sections~~
- [x] ~~**Implement plan generation workflow**~~
  - [x] ~~Loading states during API calls~~
  - [x] ~~Error handling for failed requests~~
  - [x] ~~Success feedback and plan preview~~
- [x] ~~**Add plan customization options**~~
  - [x] ~~Edit generated plans before export~~

### 🔌 Backend Integration

- [x] ~~**Set up API client with proper TypeScript types**~~
  - [x] ~~HTTP client configuration (fetch-based)~~
  - [x] ~~Request/response type definitions~~
  - [x] ~~Error handling and retry logic~~
- [x] ~~**Implement training plan generation API calls**~~
  - [x] ~~POST endpoint for plan requests~~
  - [x] ~~Proper request payload formatting~~
  - [x] ~~Response parsing and validation~~
- [x] ~~**Add API error handling and user feedback**~~
  - [x] ~~Network error recovery~~
  - [x] ~~Backend error message display~~
  - [x] ~~CORS handling via development proxy~~

### 🧪 Testing & Quality Assurance

- [x] ~~Write unit tests for core components~~
  - [x] ~~Input form validation tests~~
  - [x] ~~Plan display component tests~~
  - [x] ~~Store/state management tests~~
- [x] ~~Implement integration tests~~
  - [x] ~~End-to-end user workflows~~

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
  - Cloud build files for pr-merge to main and release

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

## 🎉 V1 Implementation Progress

### ✅ Completed Features

#### Core Architecture & Setup

- Vue 3 + TypeScript project with modern tooling
- Pinia stores for state management (training plans, settings, export)
- TypeScript type definitions aligned with backend API
- API client with CORS proxy and timeout handling
- Component architecture with layout system

#### UI/UX Implementation

- Responsive layout with AppHeader, AppFooter, AppLayout
- Advanced TrainingPlanForm with comprehensive filtering system
- TrainingPlanDisplay with animated loading states
- Clean design system without gradients, mono-color approach
- Alternating table rows for improved readability

#### Backend Integration

- Complete API client with error handling
- CORS handling via Vite proxy configuration
- Filter system matching backend Pydantic models
- Real-time form validation and submission
- Loading state management during 60s backend calls

#### User Experience

#### User Experience Features

- Sliding loading animation for plan generation
- Proper loading state prioritization (replaces old plan during generation)
- Error handling with user-friendly messages
- Responsive form with three-state filter logic

#### PDF Export System

- Server-side PDF generation using backend `/export-pdf` endpoint
- Export button integrated into training plan display
- Loading states and error handling for export process
- Type-safe request/response handling with proper TypeScript integration### 🚧 Currently Working On

- Final polish and user experience improvements

### 📋 Next Priority Items

#### V1 Polish & Enhancements

1. **Layout Improvements**: Header and footer should fill full page width
   - Change base layout to full width design
   - Add container elements for content width control
   - Maintain responsive design principles

2. **Settings UX Enhancement**: Replace help text with hoverable question mark icons
   - Add detailed tooltips for all filter settings
   - Explain difficulty levels, training types, and swimming techniques
   - Improve user understanding of advanced options

3. **Creative Assistance**: Add "I'm feeling lucky" sample prompt button
   - Integrate with backend endpoint for sample prompts
   - Help users who need inspiration for training requests
   - Backend endpoint needs implementation

4. **Pool Length Integration**: Include pool length in training plan requests
   - Add pool length field to the form
   - Integrate into prompt generation for more accurate plans
   - Consider common pool sizes (25m, 50m, yards)

5. **Intensity Legend**: Add hoverable help for intensity abbreviations
   - Question mark icon next to intensity column header
   - Tooltip explaining common swimming intensity abbreviations
   - Improve plan readability for new users

#### Core Features (Remaining)

- None! All core V1 features have been implemented ✅

### 🔄 Remaining V1 Tasks

#### Future Enhancements (V2 Scope)

- **Plan Customization**: Edit generated plans before export, save drafts locally
- **Testing Suite**: Unit tests, integration tests, accessibility improvements
- **Production Optimization**: Bundle analysis, asset optimization, deployment configuration
- **Advanced Features**: Custom filename generation, multiple export formats

### 🎯 Implementation Notes

#### Decisions Made

- **CORS Handling**: Development proxy chosen over backend CORS configuration for V1 simplicity
- **Design System**: Mono-color approach without gradients for clean, consistent appearance
- **Loading States**: Prioritized loading animation over existing plan display for better UX
- **Mobile Support**: Deferred to V2 to focus on desktop experience first

#### Known Issues for Future

- **Elite Athlete Filter**: Backend database missing plans for elite athletes (backend issue)
- **Mobile Responsiveness**: CSS media queries removed for V1, need re-implementation
- **Plan Editing**: Architecture ready but UI implementation deferred

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
