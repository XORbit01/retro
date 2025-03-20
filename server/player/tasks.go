package player

import "github.com/XORbit01/retro/shared"

func (p *Player) addTask(target string, typeTask int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Tasks[target] = shared.Task{
		Type:  typeTask,
		Error: "",
	}
}

func (p *Player) removeTask(target string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.Tasks, target)
}

func (p *Player) errorTask(target string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	task, ok := p.Tasks[target]
	if ok {
		task.Error = err.Error()
		p.Tasks[target] = task
	}
}
