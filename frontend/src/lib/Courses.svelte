<script lang="ts">
    import { FileCode, BrainCircuit, Shield, Cpu, Code, AppWindow, SquareTerminal, FileText, Search} from 'lucide-svelte'

    const vassarCourses = [
        {name: "Data Structures", icon: FileCode, courseNum: "CMPU 102"},
        {name: "Foundations of Computer science", icon: SquareTerminal, courseNum: "CMPU 145"},
        {name: "Analysis of Algorithms", icon: BrainCircuit, courseNum: "CMPU 241"},
        {name: "Theory of Computation", icon: Code, courseNum: "CMPU 240"},
        {name: "Natural Language Processing", icon: FileText, courseNum: "CMPU 366"},
        {name: "Problem Solving and abstraction", icon: Search, courseNum: "CMPU 101"},
        {name: "Compilers", icon: Cpu, courseNum: "CMPU 331"}
    ];

    const berkeleyCourses = [
        {name: "Computer Security", icon: Shield, courseNum: "CS 161"},
        {name: "Computer Architecture", icon: Cpu, courseNum: "CS 61C"},
        {name: "Software Engineering", icon: AppWindow, courseNum: "CS 169"}
    ];
</script>

<div class="coursework-subsection">
    <h2 class="coursework-title">Relevant Coursework</h2>


    <div class="schools-container">
        <div class="school-column">
            <h3 class="school-name">Vassar College</h3>
            <div class="course-grid">
                {#each vassarCourses as course}
                <div class="course-item">
                    <div class="course-icon">
                        <svelte:component this={course.icon} size={20}/>
                    </div>
                    <span class="course-name">{course.name}</span>
                    <span class="course-tooltip">{course.courseNum}</span>
                </div>
                {/each}
            </div>
        </div>

        <div class="school-column">
            <h3 class="school-name">UC Berkeley</h3>
            <div class="course-grid-berkeley">
                {#each berkeleyCourses as course}
                <div class="course-item">
                    <div class="course-icon">
                        <svelte:component this={course.icon} size={20}/>
                    </div>
                    <span class="course-name">{course.name}</span>
                    <span class="course-tooltip">{course.courseNum}</span>
                </div>
                {/each}
            </div>
        </div>

    </div>

</div>


<style>
    .coursework-subsection {
        margin-top: 60px;
        padding-top: 40px;
        border-top: 1px solid rgba(255, 255, 255, 0.1);

    }

    .coursework-title {
        font-size: 1.8rem;
        font-weight: 600;
        color: #f0f6fc;
        margin-bottom: 30px;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 12px;
    }

    .schools-container {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 40px;
        max-width: 900px;
        margin: 0 auto;
    }

    .course-grid {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 15px;
        max-width: 700px;
        margin: 0 auto;
    }

    .course-grid-berkeley {
        display: flex;
        flex-direction: column;
        gap: 12px;
        min-width: 23vh;
    }

    .school-column {
        display: flex;
        flex-direction: column;
    }

    .school-name {
        font-size: 1.3rem;
        font-weight: 600;
        color: #00d4aa;
        margin-bottom: 20px;
        text-align: center;
        padding-bottom: 10px;
        border-bottom: 2px solid rgba(0, 212, 170, 0.3);
    }

    .course-item {
        display: flex;
        align-items: center;
        gap: 12px;
        background: rgba(33, 38, 45, 0.4);
        backdrop-filter: blur(10px);
        -webkit-backdrop-filter: blur(10px);
        border: 1px solid rgba(255, 255, 255, 0.08);
        border-radius: 10px;
        padding: 12px 16px;
        transition: all 0.3s ease;
    }

    .course-item:hover {
        border-color: rgba(0, 212, 170, 0.3);
        background: rgba(33, 38, 45, 0.6);
        transform: translateX(5px);
    }

    .course-icon {
        color: #00d4aa;
        display: flex;
        align-items: center;
    }

    .course-name {
        color: #f0f6fc;
        font-size: 0.95rem;
        font-weight: 500;
    }
    @media (max-width: 768px) {
        .coursework-title {
            font-size: 1.5rem;
        }
        
        .course-grid {
            grid-template-columns: 1fr;
        }
    }

    .course-tooltip {
        position: absolute;
        top: -35px;
        left: 50%;
        transform: translateX(-50%);
        background: rgba(0, 212, 170, 0.95);
        color: #0d1117;
        padding: 6px 12px;
        border-radius: 6px;
        font-size: 0.85rem;
        font-weight: 600;
        white-space: nowrap;
        opacity: 0;
        visibility: hidden;
        transition: all 0.3s ease;
        pointer-events: none;
        box-shadow: 0 4px 12px rgba(0, 212, 170, 0.3);
    }

    .course-tooltip::after {
        content: '';
        position: absolute;
        bottom: -4px;
        left: 50%;
        transform: translateX(-50%);
        width: 0;
        height: 0;
        border-left: 5px solid transparent;
        border-right: 5px solid transparent;
        border-top: 5px solid rgba(0, 212, 170, 0.95);
    }
    .course-item:hover .course-tooltip {
        opacity: 1;
        visibility: visible;
    }
</style>