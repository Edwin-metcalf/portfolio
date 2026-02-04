<script lang="ts">
    import { FileCode, BrainCircuit, Shield, Cpu, Code, AppWindow, SquareTerminal, FileText, Search, Network} from 'lucide-svelte'

    const vassarCourses = [
        {name: "Data Structures", icon: FileCode, courseNum: "CMPU 102"},
        {name: "Foundations of Computer Science", icon: SquareTerminal, courseNum: "CMPU 145"},
        {name: "Analysis of Algorithms", icon: BrainCircuit, courseNum: "CMPU 241"},
        {name: "Theory of Computation", icon: Code, courseNum: "CMPU 240"},
        {name: "Natural Language Processing", icon: FileText, courseNum: "CMPU 366"},
        {name: "Problem Solving and abstraction", icon: Search, courseNum: "CMPU 101"},
        {name: "Compilers", icon: Cpu, courseNum: "CMPU 331"},
        {name: "Computer Networks", icon: Network, courseNum: "CMPU 375"}
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
                        <svelte:component this={course.icon} size={15}/>
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
                        <svelte:component this={course.icon} size={15}/>
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
    }

    .coursework-title {
        font-size: 3rem;
        font-weight: 600;
        color: #f0f6fc;
        text-align: center;
        margin-bottom: 50px;
        position: relative;
    }
        
    .coursework-title::after {
        content: '';
        display: block;
        width: 60px;
        height: 3px;
        background: #00d4aa;
        margin: 20px auto;
        border-radius: 2px;
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
        margin: 5 auto;
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
        display: grid;
        grid-template-columns: 28px 1fr auto;
        align-items: center;
        gap: 0px;
        background: none;
        border: none;
        border-radius: 0;
        border-bottom: 1px solid rgba(255,255,255,0.06);
        margin-bottom: 12px;
        padding: 5px 0;
    }

    .course-icon {
        color: #00d4aa;
        display: flex;
        align-items: center;
        opacity: 0.8;
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
        color: rgba(255,255,255,0.5);
        font-size: 0.85rem;
        font-family: monospace;
        display: none;
    }
</style>