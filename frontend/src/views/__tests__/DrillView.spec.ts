import en from '@/locales/en.json'
import { useDrillsStore } from '@/stores/drills'
import type { Drill } from '@/types'
import { createTestingPinia } from '@pinia/testing'
import { mount, flushPromises } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createI18n } from 'vue-i18n'
import DrillView from '../DrillView.vue'

const i18n = createI18n({
    legacy: false,
    locale: 'en',
    messages: {
        en,
    },
})

vi.mock('vue3-toastify', () => ({
    toast: {
        error: vi.fn(),
    },
}))

const pushMock = vi.fn()

vi.mock('vue-router', async (importOriginal) => {
    const actual = (await importOriginal()) as typeof import('vue-router')
    return {
        ...actual,
        useRoute: vi.fn(() => ({
            params: {
                id: 'seestern.png',
            },
        })),
        useRouter: vi.fn(() => ({
            push: pushMock,
        })),
    }
})

const mockDrill: Drill = {
    slug: 'Starfish',
    targets: ['Gliding', 'Water Resistance', 'Water Feel'],
    short_description: 'Exercise to experience water resistance and gliding ability.',
    img_name: 'seestern.png',
    img_description:
        'The swimmer pushes off the pool wall and stretches arms and legs far away from the body.',
    title: 'The Starfish - A simple exercise for water feel while gliding',
    description: [
        'Correct gliding technique is a prerequisite for learning the front crawl, but is also important for other strokes.',
        'This exercise is about getting a feel for how water resistance affects gliding distance.',
    ],
    video_url: [
        'https://www.youtube.com/watch?v=xtuRdL8PTSA&list=PLylrIIV_u33iQsIG9Dqj_B93B6KDJnsJ0&index=1',
    ],
    styles: ['General'],
    difficulty: 'Easy',
    target_groups: ['Beginner'],
}

describe('DrillView.vue', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        pushMock.mockClear()
    })

    it('fetches drill on mount', async () => {
        mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                    }),
                    i18n,
                ],
            },
        })

        const drillsStore = useDrillsStore()
        expect(drillsStore.fetchDrill).toHaveBeenCalledWith('seestern.png', 'en')
    })

    it('displays loading state', async () => {
        const wrapper = mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            drills: {
                                isLoading: true,
                                currentDrill: null,
                                error: null,
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        expect(wrapper.find('.loading-state').exists()).toBe(true)
        expect(wrapper.find('.loading-spinner').exists()).toBe(true)
        expect(wrapper.find('.drill-header').exists()).toBe(false)
    })

    it('displays error state', async () => {
        const wrapper = mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            drills: {
                                isLoading: false,
                                currentDrill: null,
                                error: 'Test Error',
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        expect(wrapper.find('.error-state').exists()).toBe(true)
        expect(wrapper.find('.error-state').text()).toContain('Test Error')
    })

    it('displays drill content when loaded', async () => {
        const wrapper = mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            drills: {
                                isLoading: false,
                                currentDrill: mockDrill,
                                error: null,
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        expect(wrapper.find('.drill-header').exists()).toBe(true)
        expect(wrapper.find('.drill-title').text()).toBe(mockDrill.title)
        expect(wrapper.find('.drill-short-description').text()).toBe(mockDrill.short_description)
    })

    it('displays difficulty tag with correct class', async () => {
        const wrapper = mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            drills: {
                                isLoading: false,
                                currentDrill: mockDrill,
                                error: null,
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        const difficultyTag = wrapper.find('.difficulty-badge')
        expect(difficultyTag.exists()).toBe(true)
        expect(difficultyTag.text()).toBe('Easy')
        expect(difficultyTag.classes()).toContain('easy')
    })

    it('displays all style tags', async () => {
        const wrapper = mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            drills: {
                                isLoading: false,
                                currentDrill: mockDrill,
                                error: null,
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        const styleTags = wrapper.findAll('.meta-tag.style')
        expect(styleTags.length).toBe(mockDrill.styles.length)
        expect(styleTags[0]?.text()).toBe('General')
    })

    it('displays description paragraphs', async () => {
        const wrapper = mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            drills: {
                                isLoading: false,
                                currentDrill: mockDrill,
                                error: null,
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        const paragraphs = wrapper.findAll('.description-text p')
        expect(paragraphs.length).toBe(mockDrill.description.length)
        expect(paragraphs[0]?.text()).toBe(mockDrill.description[0])
    })

    it('displays target tags', async () => {
        const wrapper = mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            drills: {
                                isLoading: false,
                                currentDrill: mockDrill,
                                error: null,
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        const targetTags = wrapper.findAll('.target-chip')
        expect(targetTags.length).toBe(mockDrill.targets.length)
        mockDrill.targets.forEach((target, index) => {
            expect(targetTags[index]?.text()).toContain(target)
        })
    })

    it('displays target group tags', async () => {
        const wrapper = mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            drills: {
                                isLoading: false,
                                currentDrill: mockDrill,
                                error: null,
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        const targetGroupTags = wrapper.findAll('.meta-tag.group')
        expect(targetGroupTags.length).toBe(mockDrill.target_groups.length)
        expect(targetGroupTags[0]?.text()).toBe('Beginner')
    })

    it('displays YouTube video when video_url is provided', async () => {
        const wrapper = mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            drills: {
                                isLoading: false,
                                currentDrill: mockDrill,
                                error: null,
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        const videoSection = wrapper.find('.video-section')
        expect(videoSection.exists()).toBe(true)

        const iframe = wrapper.find('.video-iframe')
        expect(iframe.exists()).toBe(true)
        expect(iframe.attributes('src')).toContain('youtube.com/embed/xtuRdL8PTSA')
    })

    it('does not display video section when video_url is empty', async () => {
        const drillWithoutVideo: Drill = {
            ...mockDrill,
            video_url: [],
        }

        const wrapper = mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            drills: {
                                isLoading: false,
                                currentDrill: drillWithoutVideo,
                                error: null,
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        const videoSection = wrapper.find('.video-section')
        expect(videoSection.exists()).toBe(false)
    })

    it('does not display video section when video_url contains empty strings', async () => {
        const drillWithEmptyVideo: Drill = {
            ...mockDrill,
            video_url: [''],
        }

        const wrapper = mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            drills: {
                                isLoading: false,
                                currentDrill: drillWithEmptyVideo,
                                error: null,
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        const videoSection = wrapper.find('.video-section')
        expect(videoSection.exists()).toBe(false)
    })

    it('generates correct image URL', async () => {
        const wrapper = mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        initialState: {
                            drills: {
                                isLoading: false,
                                currentDrill: mockDrill,
                                error: null,
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        const img = wrapper.find('.drill-image')
        expect(img.exists()).toBe(true)
        expect(img.attributes('src')).toBe('https://storage.googleapis.com/undefined/seestern.png')
        expect(img.attributes('alt')).toBe(mockDrill.img_description)
    })

    it('redirects to home when drill is not found', async () => {
        const { toast } = await import('vue3-toastify')

        mount(DrillView, {
            global: {
                plugins: [
                    createTestingPinia({
                        createSpy: vi.fn,
                        stubActions: false,
                        initialState: {
                            drills: {
                                isLoading: false,
                                currentDrill: null,
                                error: null,
                            },
                        },
                    }),
                    i18n,
                ],
            },
        })

        const drillsStore = useDrillsStore()
        // Mock fetchDrill to return null
        vi.mocked(drillsStore.fetchDrill).mockResolvedValue(null)

        await flushPromises()

        // The component should have called fetchDrill which returned null
        // This triggers noDrillFound() which shows toast and redirects
        expect(toast.error).toHaveBeenCalled()
        expect(pushMock).toHaveBeenCalledWith('/')
    })

    describe('YouTube video ID extraction', () => {
        it('extracts video ID from standard YouTube URL', async () => {
            const drillWithStandardUrl: Drill = {
                ...mockDrill,
                video_url: ['https://www.youtube.com/watch?v=abc123def45'],
            }

            const wrapper = mount(DrillView, {
                global: {
                    plugins: [
                        createTestingPinia({
                            createSpy: vi.fn,
                            initialState: {
                                drills: {
                                    isLoading: false,
                                    currentDrill: drillWithStandardUrl,
                                    error: null,
                                },
                            },
                        }),
                        i18n,
                    ],
                },
            })

            const iframe = wrapper.find('.video-iframe')
            expect(iframe.attributes('src')).toContain('youtube.com/embed/abc123def45')
        })

        it('extracts video ID from youtu.be short URL', async () => {
            const drillWithShortUrl: Drill = {
                ...mockDrill,
                video_url: ['https://youtu.be/abc123def45'],
            }

            const wrapper = mount(DrillView, {
                global: {
                    plugins: [
                        createTestingPinia({
                            createSpy: vi.fn,
                            initialState: {
                                drills: {
                                    isLoading: false,
                                    currentDrill: drillWithShortUrl,
                                    error: null,
                                },
                            },
                        }),
                        i18n,
                    ],
                },
            })

            const iframe = wrapper.find('.video-iframe')
            expect(iframe.attributes('src')).toContain('youtube.com/embed/abc123def45')
        })

        it('handles multiple videos', async () => {
            const drillWithMultipleVideos: Drill = {
                ...mockDrill,
                video_url: [
                    'https://www.youtube.com/watch?v=video111111',
                    'https://www.youtube.com/watch?v=video222222',
                ],
            }

            const wrapper = mount(DrillView, {
                global: {
                    plugins: [
                        createTestingPinia({
                            createSpy: vi.fn,
                            initialState: {
                                drills: {
                                    isLoading: false,
                                    currentDrill: drillWithMultipleVideos,
                                    error: null,
                                },
                            },
                        }),
                        i18n,
                    ],
                },
            })

            const iframes = wrapper.findAll('.video-iframe')
            expect(iframes.length).toBe(2)
            expect(iframes[0]?.attributes('src')).toContain('video111111')
            expect(iframes[1]?.attributes('src')).toContain('video222222')
        })
    })

    describe('difficulty levels', () => {
        it.each([
            ['Easy', 'easy'],
            ['Medium', 'medium'],
            ['Hard', 'hard'],
        ])('displays %s difficulty with %s class', async (difficulty, expectedClass) => {
            const drillWithDifficulty: Drill = {
                ...mockDrill,
                difficulty: difficulty as 'Easy' | 'Medium' | 'Hard',
            }

            const wrapper = mount(DrillView, {
                global: {
                    plugins: [
                        createTestingPinia({
                            createSpy: vi.fn,
                            initialState: {
                                drills: {
                                    isLoading: false,
                                    currentDrill: drillWithDifficulty,
                                    error: null,
                                },
                            },
                        }),
                        i18n,
                    ],
                },
            })

            const difficultyTag = wrapper.find('.difficulty-badge')
            expect(difficultyTag.text()).toBe(difficulty)
            expect(difficultyTag.classes()).toContain(expectedClass)
        })
    })
})
